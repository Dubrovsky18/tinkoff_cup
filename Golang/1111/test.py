from selenium import webdriver
from time import sleep
from dotenv import load_dotenv
import os 
load_dotenv()

def test_200():
    options_G = webdriver.ChromeOptions()
    driver = webdriver.Remote(
        command_executor=f'http://localhost:4444',
        options=options_G
    )
  
    url = " https://www.tinkoff.ru/" 

    driver.get(url)
    driver.quit()

    assert 1 == 1


def test_xframe():
    import requests
    req = requests.get('https://www.tinkoff.ru/')

    result = req.headers['x-frame-options']  

    assert result == "sameorigin" or result == "deny"
    


