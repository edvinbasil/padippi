# Padippi

## Using Padippi

- The following things are hardcoded and should probably be changed
    - subjects and their corresponding roles and urls
    - slots (A-H) and their corresponding subjects
    - timetable according to the slots
    - (optional) bot name and embed color
- the discord webhook url can be provided using the DISCORD_WEBHOOK_URL
  enviornment variable
- The bot is invoked using a cronjob that looks something like this

    ```cron

    00 01 * * 1-5 /usr/local/bin/classalerts daily     # send the timetabe for the day
    50 07 * * 1-5 /usr/local/bin/classalerts sendtt 8  # sends alert for the 8.00 class on 7.50
    50 08 * * 1-5 /usr/local/bin/classalerts sendtt 9
    05 10 * * 1-5 /usr/local/bin/classalerts sendtt 10
    05 11 * * 1-5 /usr/local/bin/classalerts sendtt 11
    50 12 * * 1-5 /usr/local/bin/classalerts sendtt 1
    50 13 * * 1-5 /usr/local/bin/classalerts sendtt 2
    50 14 * * 1-5 /usr/local/bin/classalerts sendtt 3
    50 15 * * 1-5 /usr/local/bin/classalerts sendtt 4
    50 16 * * 1-5 /usr/local/bin/classalerts sendtt 5

    ```

## Roadmap / improvements

- organise/refactor some of the code for clarity
- load timetable and other values from a config file, instead of hardcoding them
- add a proper http server to request the timetable for a particular day. Right
  now, this is implemented using a simple php page that calls the script using the
  correct arguments
