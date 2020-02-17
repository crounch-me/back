Feature: Create Product

  Scenario: OK
    Given I create and authenticate with a random user
    And I use this body
      """
        {
          "name": "Creation produit"
        }
      """
    When I send a "POST" request on "/products"
    Then the status code is 201
    And "$.name" is a string equal to "Creation produit"
    And "$.id" is a non empty string
