Feature: Users me

  Scenario: OK
    Given I authenticate with a random user
    When I send a "GET" request on "/me"
    Then the status code is 200
    And "$.id" is a non empty string
    And "$.email" is a non empty string
