from flask import Flask, jsonify, request
import requests
from bs4 import BeautifulSoup
import logging

app = Flask(__name__)

def scrape_property_site(suburb):
    url = f'https://www.trademe.co.nz/a/property/residential/sale/canterbury/christchurch-city/{suburb}'
    response = requests.get(url)
    soup = BeautifulSoup(response.content, 'html.parser')
    listings = []

    logging.info(f"Scraping URL: {url}")
    logging.info(f"Response status code: {response.status_code}")
    logging.info(f"Response content: {response.content[:500]}")  # Print the first 500 characters of the response

    # Check if the elements we're targeting exist in the HTML
    if not soup.find_all('div', class_='tm-property-search-card__content'):
        logging.warning("No listings found in the HTML content.")
    else:
        logging.info("Listings found in the HTML content.")

    for listing in soup.find_all('div', class_='tm-property-search-card__content'):
        try:
            title = listing.find('h2').text.strip()
            rooms = listing.find('span', class_='tm-property-search-card__features').text.strip()
            suburb = listing.find('span', class_='tm-property-search-card__location').text.strip()
            price = listing.find('span', class_='tm-property-search-card__price').text.strip()
            link = listing.find('a', class_='tm-property-search-card__link')['href']
            listings.append({'Title': title, 'Rooms': rooms, 'Suburb': suburb, 'Price': price, 'Link': link})
            logging.info(f"Processed listing: {listings[-1]}")
        except Exception as e:
            logging.error(f"Error processing listing: {e}")

    logging.info(f"Found {len(listings)} listings")
    return listings

@app.route('/api/houses', methods=['GET'])
def get_houses():
    suburb = request.args.get('suburb')
    rooms = request.args.get('rooms')
    price = request.args.get('price')

    houses = scrape_property_site(suburb)
    filtered_houses = []

    for house in houses:
        if rooms and int(house['Rooms']) != int(rooms):
            continue
        if price and int(house['Price'].replace(',', '').replace('$', '')) > int(price):
            continue
        filtered_houses.append(house)

    return jsonify(filtered_houses)

if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    app.run(port=5000)