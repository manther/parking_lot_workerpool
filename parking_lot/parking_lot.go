package parkinglot

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// / Parking lot structure and behaviors
type Parkinglot struct {
	totalComp int
	totalReg  int
	compact   []spot
	regular   []spot
}

// Constructor populates two collections of spots, and a total count of each
func NewParparkinglot(compactCount int, regularCount int) Parkinglot {
	cp := make([]spot, compactCount)
	for i := 0; i < compactCount; i++ {
		cp[i] = &compactSpot{
			occupied: false,
			lotSize:  0,
			id:       i,
		}
	}

	rg := make([]spot, regularCount)
	for i := 0; i < regularCount; i++ {
		rg[i] = &regularSpot{
			occupied: false,
			lotSize:  0,
			id:       i,
		}
	}

	return Parkinglot{
		totalComp: compactCount,
		totalReg:  regularCount,
		compact:   cp,
		regular:   rg,
	}
}

// / Parking lot getters
func (p Parkinglot) Full() bool {
	return len(p.compact)+len(p.regular) <= 0
}

func (p Parkinglot) Total() int {
	return p.CompactAvailable() + p.RegularAvailable()
}

func (p Parkinglot) CompactTot() int {
	return p.totalComp
}

func (p Parkinglot) RegularTot() int {
	return p.totalReg
}

func (p Parkinglot) CompactAvailable() int {
	return len(p.compact)
}

func (p Parkinglot) RegularAvailable() int {
	return len(p.regular)
}

// SpotType Enum
type spotType int
const (
	compact spotType = iota
	regular
)

// Spot expected behaviors
type spot interface {
	isOccupied() bool
	size() spotType
	occupy() error
	vacate()
	spotID() int
}

// Concreate spot types
type compactSpot struct {
	occupied bool
	lotSize  spotType
	id       int
}

func (s compactSpot) isOccupied() bool {
	return s.occupied
}

func (s compactSpot) size() spotType {
	return compact
}

func (s *compactSpot) occupy() error {
	if s.occupied {
		return errors.New("cannot occupy. Spot already occupied")
	}
	s.occupied = true
	return nil
}

func (s *compactSpot) vacate() {
	s.occupied = false
}

func (s compactSpot) spotID() int {
	return s.id
}

type regularSpot struct {
	occupied bool
	lotSize  spotType
	id       int
}

func (s regularSpot) isOccupied() bool {
	return s.occupied
}

func (s regularSpot) size() spotType {
	return regular
}

func (s *regularSpot) occupy() error {
	if s.occupied {
		return errors.New("cannot occupy. Spot already occupied")
	}
	s.occupied = true
	return nil
}

func (s *regularSpot) vacate() {
	s.occupied = false
}

func (s regularSpot) spotID() int {
	return s.id
}




// / Parkinlot actors, structures, and behaviors
// Parking management
// Other manager could be defined in its place
type Manager interface {
	manageGroupParkComp(numToPark int) []error
	manageSingleParkComp() error
	manageGroupParkReg(numToPark int) []error
	manageSingleParkReg() error
}

type LotSupervisor struct {
	jobs    chan parkingJob
	results chan error
	pl      *Parkinglot
	vValet  parker
}

// Worker pool
func (s LotSupervisor) worker(id int, jobs <-chan parkingJob, results chan<- error) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j.id)
		time.Sleep(time.Second)
		err := j.f(j.sp)
		fmt.Println("worker", id, "finished job", j.id)
		results <- err
	}
}

// Constructor - Supervisor
func NewLotSupervisor(pl *Parkinglot, valetCount int) LotSupervisor {
	jobs := make(chan parkingJob, valetCount*2)
	results := make(chan error, valetCount*2)
	mutex := &sync.Mutex{}

	sup := LotSupervisor{
		jobs:    jobs,
		results: results,
		pl:      pl,
		vValet:  valet{mutex: mutex},
	}

	// spin up valets
	// Valet count supplied by user
	for w := 1; w <= valetCount; w++ {
		go sup.worker(w, jobs, results)
	}
	return sup
}

// Parking jobs that are sent to the worker pool
type parkingJob struct {
	id uuid.UUID
	f  func(*[]spot) error
	sp *[]spot
}

// Supported management behaviors, multi and single for compact spots
func (s LotSupervisor) ManageGroupParkComp(numToPark int) []error {
	errs := []error{}
	wg := sync.WaitGroup{}
	for i := 0; i < numToPark; i++ {
		if len(s.pl.compact) > numToPark {
			wg.Add(1)
			go func() {
				defer wg.Done()
				id := uuid.New()
				s.jobs <- parkingJob{
					id,
					s.vValet.park,
					&s.pl.compact,
				}
				errs = append(errs, <-s.results)
			}()
		} else {
			errs = append(errs, errors.New("not enough compact spots"))
			break
		}
	}

	wg.Wait()
	return errs
}

func (s LotSupervisor) ManageSingleParkComp() error {
	if len(s.pl.compact) > 0 {
		go func() {
			s.jobs <- parkingJob{
				uuid.New(),
				s.vValet.park,
				&s.pl.compact,
			}
		}()
	} else {
		return errors.New("not enough compact spots")
	}
	res := <-s.results
	return res
}

// Supported management behaviors, multi and single for regular spots
func (s LotSupervisor) ManageGroupParkReg(numToPark int) []error {
	errs := []error{}
	wg := sync.WaitGroup{}
	for i := 0; i < numToPark; i++ {
		if len(s.pl.regular) > numToPark {
			wg.Add(1)
			go func() {
				defer wg.Done()
				s.jobs <- parkingJob{
					uuid.New(),
					s.vValet.park,
					&s.pl.regular,
				}
				errs = append(errs, <-s.results)
			}()
		} else {
			errs = append(errs, errors.New("not enough regular spots"))
			break
		}
	}

	wg.Wait()
	return errs
}

func (s LotSupervisor) ManageSingleParkReg() error {
	if len(s.pl.regular) > 0 {
		go func() {
			s.jobs <- parkingJob{
				uuid.New(),
				s.vValet.park,
				&s.pl.regular,
			}
		}()
	} else {
		return errors.New("not enough compact spots")
	}
	res := <-s.results
	return res
}

// Parking behaviors interface
type parker interface {
	park(*[]spot) error
}

// Concrete parker
type valet struct {
	mutex *sync.Mutex
}

// Park implementation with mutual exclusion on lot resources
func (p valet) park(sps *[]spot) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !(*sps)[0].isOccupied() {
		(*sps)[0].occupy()
		*sps = (*sps)[1:]

		return nil
	} else {
		return errors.New("spot is occupied but we had no record of it. Remove from the pool")
	}
}