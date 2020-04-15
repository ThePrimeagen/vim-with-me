package scheduled

import (
    "time"
)

type Scheduled struct {
    Expires int64
    Done chan interface{}
}

func (s *Scheduled) AddSeconds(seconds int64) {

    now := GetCurrentMillis()
    diff := s.Expires - now
    expired := diff < 0

    if expired {
        diff = 0
    }

    s.Expires = now + diff + seconds * 1000

    if expired {
        go func() {
            for {
                now = GetCurrentMillis()

                remaining := (s.Expires - now) / 1000
                if remaining <= 0 {
                    break
                }

                t := time.NewTimer(time.Duration(remaining) * time.Second)
                <-t.C
            }

            s.Done <- 1
        }()
    }
}

func GetCurrentMillis() int64 {
    now := time.Now()
    unixNano := now.UnixNano()
    return unixNano / 1000000
}

func NewScheduled(expireAmount int64) Scheduled {
    now := GetCurrentMillis()

    scheduled := Scheduled{
        now + expireAmount,
        make(chan interface{}),
    }

    return scheduled
}

