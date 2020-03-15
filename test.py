from selenium import webdriver
from selenium.webdriver.common.desired_capabilities import DesiredCapabilities
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support.expected_conditions import presence_of_element_located,visibility_of,visibility_of_element_located,visibility_of_all_elements_located
from  time import sleep
import sys
#This example requires Selenium WebDriver 3.13 or newer
option = webdriver.ChromeOptions()
option.add_argument('headless')
option.add_argument("--enable-logging --v=1")
d = DesiredCapabilities.CHROME
d['goog:loggingPrefs'] = { 'browser':'ALL' }
with webdriver.Chrome(options=option,desired_capabilities=d) as driver:
    wait = WebDriverWait(driver, 10)
    driver.get("https://www.google.com/maps")
    _input = driver.find_elements_by_css_selector("#searchboxinput")[0]
    _input.send_keys(sys.argv[1])
    _input.send_keys(Keys.ENTER)
    img = wait.until(visibility_of_element_located((By.CSS_SELECTOR, ".section-hero-header-image-hero > img")))
    img.click()
    sleep(1)
    driver.execute_script("console.log(document.location.href)")
    log = driver.get_log('browser')
    print(str(log[0]["message"]).split(" ")[2])

