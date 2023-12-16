package tmux

import "errors"

var (
	// ErrNilRunner is returned when a nil [Runner] argument is passed.
	ErrNilRunner = errors.New("runner is nil")
	// ErrNilSession is returned when a nil [Session] argument is passed.
	ErrNilSession = errors.New("session is nil")
	// ErrNilWindow is returned when a nil [Window] argument is passed.
	ErrNilWindow = errors.New("window is nil")
	// ErrNilPane is returned when a nil [Pane] argument is passed.
	ErrSessionClosed = errors.New("session is closed")
	// ErrSessionNotApplied is returned when an unapplied [Session] is used.
	ErrSessionNotApplied = errors.New("session is not applied")
	// ErrWindowNotApplied is returned when an unapplied [Window] is used.
	ErrWindowNotApplied = errors.New("window is not applied")
	// ErrPaneNotApplied is returned when an unapplied [Pane] is used.
	ErrPaneNotApplied = errors.New("pane is not applied")
)
