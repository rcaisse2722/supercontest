# Ross Caisse 2023-01-23
# As of January 2023, NFL regular season picks from completed 2023 season are accessible
# via https://www.superbook.com/supercontest. The picks/card is returned directly as part
# of the raw HTML, so now API calls/WebSocket communication is necessary to retrieve this 
# data. This script is used for getting all picks for each week, as well as the card for
# each week, and write this data to raw HTML files for later processing.
import os
import time
from urllib.request import Request, urlopen 

START_WEEK = 11
END_WEEK = 18
OUTPUT_DIR = "raw_data"
TOKEN = "__TOKEN__"
PICKS_BASE_URL = f"https://www.superbook.com/supercontest/2022-supercontest-week-{TOKEN}-selections"
CARD_BASE_URL = f"https://www.superbook.com/supercontest/weekly-card-{TOKEN}/"


def get_html_data(url):
    req = Request(url)
    req.add_header("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:108.0) Gecko/20100101 Firefox/108.0")
    req.add_header("Accept", "*/*")
    req.add_header("Connection", "keep-alive")
    return urlopen(req).read().decode('utf-8')

# Check if output directory exists, if not create it
folder_exists = os.path.isdir(OUTPUT_DIR)
if not folder_exists:
    print(f"Output folder {OUTPUT_DIR} does not exists. Creeating it.")
    os.makedirs(OUTPUT_DIR)

for week_num in range(START_WEEK, END_WEEK + 1):
    # Get the weekly picks raw HTML    
    weekly_pick_raw_html = get_html_data(PICKS_BASE_URL.replace(TOKEN, str(week_num)))
    # # Output picks HTML to file
    with open(f"{OUTPUT_DIR}/week_{week_num}_picks.html", "w") as f:
        f.write(weekly_pick_raw_html)

    # Get the weekly card raw HTML
    weekly_card_html = get_html_data(CARD_BASE_URL.replace(TOKEN, str(week_num).zfill(2)))
    # Output picks HTML to file
    with open(f"{OUTPUT_DIR}/week_{week_num}_card.html", "w") as f:
        f.write(weekly_card_html)

    print(f"Done writing Week {week_num} data")

    # Wait for 5 seconds for the next requests
    time.sleep(15)

