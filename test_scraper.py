from house_scraper import scrape_property_site

def main():
    suburb = "edgeware"  # Default suburb for testing
    listings = scrape_property_site(suburb)
    for listing in listings:
        print(listing)

if __name__ == '__main__':
    main()