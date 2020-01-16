package throttler

import (
	"testing"
)

type task struct {
}

func (t task) Run() {
	return
}
func TestThrottler_AddTask(t1 *testing.T) {
	type fields struct {
		tasks     []Runnable
		toRun     int
		isProcess bool
	}
	type args struct {
		task task
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name:   "Add one task",
			fields: fields{tasks: make([]Runnable, 0)},
			args:   args{task: task{}},
			want:   1,
		},
		{
			name:   "Add one task where toRun > 0",
			fields: fields{toRun: 10, tasks: make([]Runnable, 0)},
			args:   args{task: task{}},
			want:   1,
		},
		{
			name:   "Add one task where len tasks > 0",
			fields: fields{tasks: make([]Runnable, 10, 11)},
			args:   args{task: task{}},
			want:   11,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Throttler{
				tasks:     tt.fields.tasks,
				toRun:     tt.fields.toRun,
				isProcess: tt.fields.isProcess,
			}
			if got := t.AddTask(tt.args.task); got != tt.want {
				t1.Errorf("AddTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestThrottler_QueueLen(t1 *testing.T) {
	type fields struct {
		tasks     []Runnable
		toRun     int
		isProcess bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Empty queue",
			fields: fields{
				tasks:     make([]Runnable, 0, 0),
				toRun:     0,
				isProcess: false,
			},
			want: 0,
		},
		{
			name: "Not empty queue",
			fields: fields{
				tasks:     make([]Runnable, 10, 10),
				toRun:     0,
				isProcess: false,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := Throttler{
				tasks:     tt.fields.tasks,
				toRun:     tt.fields.toRun,
				isProcess: tt.fields.isProcess,
			}
			if got := t.QueueLen(); got != tt.want {
				t1.Errorf("QueueLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestThrottler_Run(t1 *testing.T) {
	tasks10 := make([]task, 10)
	for i := 0; i < 10; i++ {
		tasks10[i] = task{}
	}
	task1 := make([]task, 1)
	task1[0] = task{}

	type fields struct {
		tasks     []task
		toRun     int
		isProcess bool
	}
	type args struct {
		count int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantToRun int
		wantErr   bool
	}{
		{
			name: "empty queue",
			fields: fields{
				tasks:     make([]task, 0),
				toRun:     0,
				isProcess: false,
			},
			args:      args{10},
			wantToRun: 10,
			wantErr:   false,
		},
		{
			name: "one task in queue, full run",
			fields: fields{
				tasks:     task1,
				toRun:     0,
				isProcess: false,
			},
			args:      args{1},
			wantToRun: 0,
			wantErr:   false,
		},
		{
			name: "ten task in queue, part run",
			fields: fields{
				tasks:     tasks10,
				toRun:     0,
				isProcess: false,
			},
			args:      args{5},
			wantToRun: 0,
			wantErr:   false,
		},
		{
			name: "ten task in queue, over run",
			fields: fields{
				tasks:     tasks10,
				toRun:     0,
				isProcess: false,
			},
			args:      args{15},
			wantToRun: 5,
			wantErr:   false,
		},
		{
			name: "ten task in queue, full run with preset",
			fields: fields{
				tasks:     tasks10,
				toRun:     5,
				isProcess: false,
			},
			args:      args{5},
			wantToRun: 0,
			wantErr:   false,
		},
		{
			name: "catch the error",
			fields: fields{
				tasks:     tasks10,
				toRun:     0,
				isProcess: false,
			},
			args:      args{-5},
			wantToRun: 0,
			wantErr:   true,
		},
		{
			name: "Run with 0",
			fields: fields{
				tasks:     tasks10,
				toRun:     5,
				isProcess: false,
			},
			args:      args{0},
			wantToRun: 0,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Throttler{
				toRun:     tt.fields.toRun,
				isProcess: tt.fields.isProcess,
			}
			for _, task := range tt.fields.tasks {
				t.AddTask(task)
			}
			gotToRun, err := t.Run(tt.args.count)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToRun != tt.wantToRun {
				t1.Errorf("Run() gotToRun = %v, want %v", gotToRun, tt.wantToRun)
			}
		})
	}
}
