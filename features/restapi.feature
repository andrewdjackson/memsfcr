
Feature: MemsFCR Rest API
  # Enter feature description here

  Scenario Outline: Connect and Initialise the ECU
    Given the serial port "<port>"
    When the ConnectAndInitialise Rest API is called
    Then the ECU connection is "<connected>"
    And the ECU has been initialised "<initialised>"
    Then disconnect the ECU

  Examples:
    | port | connected | initialised |
    | /Users/andrew.jacksonglobalsign.com/ttyecu | True     | True        |
    | /dev/null | False     | False        |