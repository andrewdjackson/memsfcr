from behave import *
from hamcrest import assert_that, equal_to
import requests
import json

@given(u'the serial port "{port}"')
def step_impl(context, port):
    context.port = port


@when(u'the ConnectAndInitialise Rest API is called')
def step_impl(context):
    url = "http://localhost:8081/rosco/connect"

    r = requests.post(url, json={"Port": f"{context.port}"})
    context.response = r.json()


@then(u'the ECU connection is "{connected}"')
def step_impl(context, connected):
    value = context.response["Connected"]
    assert_that(str(value), equal_to(connected))


@then(u'the ECU has been initialised "{initialised}"')
def step_impl(context, initialised):
    value = context.response["Initialised"]
    assert_that(str(value), equal_to(initialised))

@then(u'disconnect the ECU')
def step_impl(context):
    url = "http://localhost:8081/rosco/disconnect"
    r = requests.post(url, json={"Port": f"{context.port}"})
