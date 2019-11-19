@current
Feature: Health

  Scenario: CORS
    When I send a "OPTIONS" request on "/health"
    Then the status code is 200
    And the header "Access-Control-Allow-Methods" equals "GET,OPTIONS"
    And the header "Access-Control-Allow-Origin" equals "*"
    And the header "Access-Control-Allow-Headers" equals "Content-Type,Authorization,Accept"

  Scenario: Read health
    When I send a "GET" request on "/health"
    Then the status code is 200
    And "$.alive" has bool value "true"
