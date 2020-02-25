Feature: Add a product to a list

  Scenario: OK
    Given I create and authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    And I create these products
      | name                |
      | Mon premier produit |
    And I use this body
      """
        {
          "product_id": "{{ .ProductID }}",
          "list_id": "{{ .ListID }}"
        }
      """
    When I send a "POST" request on "/lists/{{ .ListID }}/products/{{ .ProductID}}"
    Then the status code is 201
    And "$.product_id" is a string equal to "{{ .ProductID }}"
    And "$.list_id" is a string equal to "{{ .ListID }}"
