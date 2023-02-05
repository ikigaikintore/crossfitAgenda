# Crossfit Agenda

Connect crossfit using a picture of their schedule and set it in your Google Calendar

## DONE

- Download the picture from a URL resource
- Use an OCR library to get the texts
- Create a CLI to set up the days to book the days
- Register the days with the time and the exercise in your Google Calendar (or any Calendar service)
- Retry retrieve credentials if token has expired
- Handle the authorize error in Calendar and retry credentials
- Create external API HTTP or GRPC

## TODO

- ~~[Optional] try to book the dates in the app~~
- Cache the image and use it
- Cache the ocr result
- If an event is unbooked, delete it from the calendar
- ~~Use events~~

## Run it

Create a folder called `env` and put the credentials file named `crossfitagenda.json`

```bash
go build -o crossfit cmd/main.go

GOOGLE_APPLICATION_CREDENTIALS=env/crossfitagenda.json ./crossfit
```