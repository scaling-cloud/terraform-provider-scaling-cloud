# A Working Hours set is a named, reusable timezone + weekly windows. The
# timezone lives on the set and is never inferred from a targeted schedule.
# Escalation policy step conditions reference a set to drive follow-the-sun
# routing ("page the UK rotation during UK hours, the US rotation otherwise").
resource "scaling_working_hours" "uk_office" {
  name     = "UK office hours"
  timezone = "Europe/London"

  window {
    days  = [1, 2, 3, 4, 5] # Monday–Friday (ISO weekday numbers, 1=Mon)
    start = "09:00"
    end   = "17:00"
  }
}

resource "scaling_working_hours" "us_office" {
  name     = "US office hours"
  timezone = "America/New_York"

  window {
    days  = [1, 2, 3, 4, 5]
    start = "08:00"
    end   = "18:00"
  }
}
