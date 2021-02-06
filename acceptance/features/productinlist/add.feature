@product-in-list
Feature: Add a product to a list

  Scenario: OK - User product
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    And I create these products
      | name                |
      | Mon premier produit |
    And I use this body
      """
        {
          "productID": "{{ .ProductID }}",
          "listID": "{{ .ListID }}"
        }
      """
    When I send a "POST" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}"
    Then the status code is 201

  # Scenario: OK - Default product
  #   Given I authenticate with a random user
  #   And I create these lists
  #     | name                           |
  #     | Récupération listes de courses |
  #   And I use this body
  #     """
  #       {
  #         "productId": "a40a3f16-ae0d-4b2a-884d-c8a08bb13aa4",
  #         "listId": "{{ .ListID }}"
  #       }
  #     """
  #   When I send a "POST" request on "/listing/lists/{{ .ListID }}/products/a40a3f16-ae0d-4b2a-884d-c8a08bb13aa4"
  #   Then the status code is 201

  Scenario: KO - List doesn't belong to user
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    And I authenticate with a random user
    And I create these products
      | name                |
      | Mon premier produit |
    And I use this body
      """
        {
          "productId": "{{ .ProductID }}",
          "listId": "{{ .ListID }}"
        }
      """
    When I send a "POST" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}"
    Then the status code is 403
    And "$.error" has string value "forbidden-error"

  # Scenario: KO - Product doesn't belong to user
  #   Given I authenticate with a random user
  #   And I create these products
  #     | name                |
  #     | Mon premier produit |
  #   And I authenticate with a random user
  #   And I create these lists
  #     | name                           |
  #     | Récupération listes de courses |
  #   And I use this body
  #     """
  #       {
  #         "productId": "{{ .ProductID }}",
  #         "listId": "{{ .ListID }}"
  #       }
  #     """
  #   When I send a "POST" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}"
  #   Then the status code is 403
  #   And "$.error" has string value "forbidden-error"

  Scenario: KO - Product already in list
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    And I create these products
      | name                |
      | Mon premier produit |
    And I use this body
      """
        {
          "productID": "{{ .ProductID }}",
          "listID": "{{ .ListID }}"
        }
      """
    And I send a "POST" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID}}"
    When I send a "POST" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID}}"
    Then the status code is 409
    And "$.error" has string value "product-already-in-list-error"

  Scenario: KO - List id is not an UUID
    Given I authenticate with a random user
    And I use this body
      """
        {
          "productID": "{{ .ProductID }}",
          "listID": "a"
        }
      """
    When I send a "POST" request on "/listing/lists/{{ .ListID }}/products/00000000-0000-0000-0000-000000000000"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    And "$.fields[0].name" has string value "listID"
    And "$.fields[0].error" has string value "uuid"

  Scenario: KO - Product id is not an UUID
    Given I authenticate with a random user
    And I use this body
      """
        {
          "productID": "a",
          "listID": "00000000-0000-0000-0000-000000000000"
        }
      """
    When I send a "POST" request on "/listing/lists/00000000-0000-0000-0000-000000000000/products/a"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    And "$.fields[0].name" has string value "productID"
    And "$.fields[0].error" has string value "uuid"

  Scenario: KO - List not found
    Given I authenticate with a random user
    And I use this body
      """
        {
          "productID": "00000000-0000-0000-0000-000000000000",
          "listID": "00000000-0000-0000-0000-000000000000"
        }
      """
    When I send a "POST" request on "/listing/lists/00000000-0000-0000-0000-000000000000/products/00000000-0000-0000-0000-000000000000"
    Then the status code is 404
    And "$.error" has string value "list-not-found-error"

  Scenario: KO - Product not found
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    And I use this body
      """
        {
          "productID": "00000000-0000-0000-0000-000000000000",
          "listID": "{{ .ListID }}"
        }
      """
    When I send a "POST" request on "/listing/lists/{{ .ListID }}/products/00000000-0000-0000-0000-000000000000"
    Then the status code is 404
    And "$.error" has string value "product-not-found-error"
