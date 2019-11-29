Feature: List

  Scenario: CORS
    When I send a "OPTIONS" request on "/lists"
    Then the status code is 200
    And the header "Access-Control-Allow-Methods" equals "POST,GET,OPTIONS"
    And the header "Access-Control-Allow-Origin" equals "*"
    And the header "Access-Control-Allow-Headers" equals "Content-Type,Authorization,Accept"

  Scenario: Create a list
    Given I create and authenticate with a random user
    And I use this body
      """
        {
          "name": "Creation liste de courses"
        }
      """
    When I send a "POST" request on "/lists"
    Then the status code is 201
    And "$.name" is a string equal to "Creation liste de courses"
    And "$.id" is a non empty string

  Scenario: Get owner's lists
    Given I create and authenticate with a random user
    And I create this lists
      | name                           |
      | Récupération listes de courses |
    When I send a "GET" request on "/lists"
    Then the status code is 200
    And "$.name[0]" is a string equal to "Récupération listes de courses"
    And "$.id[0]" is a non empty string
