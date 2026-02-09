# RSVP

RSVP is a system to manage invitations.

## Features

- Create events
- Upload CSV files with guest information
- Search for guests
- Confirm attendance
- Export CSV files

## Usage

### Create an event

1. Go to [Upload](/upload)
2. Upload a CSV file with the following columns:
   - `nombre`: name
   - `apellido`: last name
   - `code`: code is an arbitrary string to identify a invitation or group of invitations example: "family_Doe" or "lksdgj123"
   - `phone`: phone is unique for each invitation and code combination
   - `evento`: event is a identifier of the event alow filter by event on path /rsvp/:event if it is empty create a default value "default"
   - `respuesta`: response (optional) if it is empty create a default value "sin respuesta"

### Confirm attendance

1. Go to [RSVP](/rsvp) for default event or your custom event RSVP(event_code): [rsvp/event_code](/rsvp/event_code)
2. Search for your invitation code or name
3. Select an invitation or invitations (if you're going with more than one person with the same code)
4. Confirm attendance
5. Enjoy the event

### Check an event

1. Go to [Upload](/upload) there you can check progress of any event

### Export the CSV file

1. Go to [Export](/export) for download "default" event
2. Go to [Export](/export/event_code) for download "event_code" event   