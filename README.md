# Crossfit Agenda

Connect crossfit using a picture of their schedule and set it in your Google Calendar

## DONE

- Download the picture from a URL resource
- Use an OCR library to get the texts
- Create a CLI to set up the days to book the days
- Register the days with the time and the exercise in your Google Calendar (or any Calendar service)

## TODO

- Handle the authorize error in Calendar and retry credentials
- ~~[Optional] try to book the dates in the app~~
- Cache the image and use it
- Cache the ocr result
- Retry retrieve credentials if token has expired
- Use events