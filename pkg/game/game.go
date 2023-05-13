package game

import (
	"time"

	"github.com/name212/mines/pkg/utils"
)

type Status int

const (
	New     Status = 0
	Error   Status = 1
	Started Status = 2
	Win     Status = 3
	Lose    Status = 4
)

type Mines struct {
	field            *Field
	outputFieldCache *Field
	bombs            int
	bombsMarked      int
	status           Status
	startedAt        time.Time
	finishedAt       time.Time
	finishChecker    *finishGameChecker
	bombsGenerator   RandNumGenerator
}

func NewGame(width, height int, bombs int, bombsGenerator RandNumGenerator) *Mines {
	field := newEmptyField(width, height)

	return &Mines{
		bombs:          bombs,
		field:          field,
		status:         New,
		finishChecker:  &finishGameChecker{},
		bombsGenerator: bombsGenerator,
	}
}

func (m *Mines) Open(x, y int) {
	if m.status == New {
		err := initField(m.field, m.bombs, m.bombsGenerator, x, y)
		if err != nil {
			m.status = Error
			return
		}
		m.startedAt = time.Now()
		m.calcOpened(x, y)
		m.status = Started
		m.outputFieldCache = nil
		return
	}

	if m.field.Cell(x, y).MarkedAsBomb() {
		return
	}

	m.outputFieldCache = nil

	m.calcOpened(x, y)
	m.calcFinished()
}

func (m *Mines) SwitchMarkAsBomb(x, y int) {
	c := m.field.Cell(x, y)
	c.markedAsBomb = !c.MarkedAsBomb()
	if c.MarkedAsBomb() {
		m.bombsMarked++
	} else {
		m.bombsMarked--
	}

	m.outputFieldCache = nil
}

func (m *Mines) calcOpened(x, y int) {
	cell := m.field.Cell(x, y)
	cell.opened = true
	if cell.HasBomb() {
		m.lose()
		return
	}

	emptyCellsQueue := utils.NewQueue[Cell]()
	emptyCellsQueue.Add(cell)

	cellsCalculator := newOpenCellChecker()
	openCellWalkerInst := newOpenCellWalker(emptyCellsQueue)
	visitedEmptyCells := make(map[int]struct{})

	for !emptyCellsQueue.IsEmpty() {
		c := emptyCellsQueue.Dequeue()
		pos := m.field.linearFromPos(c.X(), c.Y())
		if _, ok := visitedEmptyCells[pos]; ok {
			continue
		}
		if c.BombsAround() == 0 {
			visitedEmptyCells[pos] = struct{}{}
		}

		if c.BombsAround() > 0 {
			cellsCalculator.reset()
			m.field.WalkAroundCell(c.X(), c.Y(), cellsCalculator)
			if cellsCalculator.hasIncorrectBombMark {
				m.lose()
				return
			}

			if cellsCalculator.markedAsBomb != c.BombsAround() {
				continue
			}
		}

		m.field.WalkAroundCell(c.X(), c.Y(), openCellWalkerInst)
	}
}

func (m *Mines) lose() {
	m.status = Lose
}

func (m *Mines) calcFinished() {
	if m.status != Started {
		return
	}

	m.finishChecker.reset()
	m.field.Walk(m.finishChecker)

	if m.finishChecker.bombOpened {
		m.status = Lose
		return
	}

	shouldOpened := m.field.Size() - m.bombs
	if shouldOpened == m.finishChecker.openedCells {
		m.finishedAt = time.Now()
		m.status = Win
	}
}

func (m *Mines) Field() *Field {
	if m.outputFieldCache == nil {
		m.outputFieldCache = m.field.Clone()
	}

	return m.outputFieldCache
}

func (m *Mines) Status() Status {
	return m.status
}

func (m *Mines) StartedAt() time.Time {
	return m.startedAt
}

func (m *Mines) FinishedAt() time.Time {
	return m.finishedAt
}

func (m *Mines) Bombs() int {
	return m.bombs
}

func (m *Mines) BombsMarked() int {
	return m.bombsMarked
}

type finishGameChecker struct {
	openedCells int
	bombOpened  bool
}

func (c *finishGameChecker) HandleCell(cell *Cell) (stop bool) {
	if !cell.Opened() {
		return !cell.MarkedAsBomb()
	}

	if cell.HasBomb() {
		c.bombOpened = true
		return true
	}

	c.openedCells++
	return false
}

func (c *finishGameChecker) reset() {
	c.openedCells = 0
	c.bombOpened = false
}

type openCellChecker struct {
	hasIncorrectBombMark bool
	markedAsBomb         int
}

func newOpenCellChecker() *openCellChecker {
	return &openCellChecker{}
}

func (c *openCellChecker) HandleCell(cell *Cell) (stop bool) {
	if cell.MarkedAsBomb() {
		if cell.HasBomb() {
			c.markedAsBomb++
		} else {
			c.hasIncorrectBombMark = true
			return true
		}
	} else {
		if cell.HasBomb() {
			c.markedAsBomb = 0
			return true
		}
	}

	return false
}

func (c *openCellChecker) reset() {
	c.hasIncorrectBombMark = false
	c.markedAsBomb = 0
}

type openCellWalker struct {
	emptyCellQueue *utils.Queue[Cell]
}

func newOpenCellWalker(emptyCellQueue *utils.Queue[Cell]) *openCellWalker {
	return &openCellWalker{
		emptyCellQueue: emptyCellQueue,
	}
}

func (w *openCellWalker) HandleCell(c *Cell) bool {
	if !c.Opened() && c.BombsAround() == 0 && !c.HasBomb() {
		w.emptyCellQueue.Add(c)
	}

	if !c.HasBomb() {
		c.opened = true
	}

	return false
}
