Feature: Search Default Products by Name

  Scenario: OK - Beginning of the name
    Given I create and authenticate with a random user
    And I use this body
      """
        {
          "name": "len"
        }
      """
    When I send a "POST" request on "/products/search"
    Then the status code is 200
    And "$[0].id" is a non empty string
    And "$[0].name" is a string equal to "Lentille"
    And "$[0].category.name" is a string equal to "Epicerie"

  Scenario: OK - Middle of the name
    Given I create and authenticate with a random user
    And I use this body
      """
        {
          "name": "mor"
        }
      """
    When I send a "POST" request on "/products/search"
    Then the status code is 200
    And "$[0].id" is a non empty string
    And "$[0].name" is a string equal to "Saucisse de Morteau"
    And "$[0].category.name" is a string equal to "Boucherie"
