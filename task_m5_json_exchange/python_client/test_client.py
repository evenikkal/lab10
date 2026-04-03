"""Tests for M5 Python client using mocked HTTP responses."""

import json
import unittest
from unittest.mock import patch, MagicMock

from client import (
    Address, OrderItem, Order,
    create_order, get_order, list_orders,
)


SAMPLE_RESPONSE = {
    "id": 1,
    "customer_id": 101,
    "items": [
        {"product_id": 1, "product_name": "Laptop Pro", "quantity": 1, "unit_price": 1499.99},
        {"product_id": 2, "product_name": "USB Hub",    "quantity": 3, "unit_price": 29.99},
    ],
    "ship_to": {
        "street": "Lenina 5, apt 12",
        "city": "Saint Petersburg",
        "country": "Russia",
        "zip": "190000",
    },
    "total_amount": 1589.96,
    "status": "pending",
    "created_at": "2024-01-01T00:00:00Z",
}


def make_sample_order() -> Order:
    return Order(
        customer_id=101,
        items=[
            OrderItem(product_id=1, product_name="Laptop Pro", quantity=1, unit_price=1499.99),
            OrderItem(product_id=2, product_name="USB Hub",    quantity=3, unit_price=29.99),
        ],
        ship_to=Address(
            street="Lenina 5, apt 12",
            city="Saint Petersburg",
            country="Russia",
            zip="190000",
        ),
    )


class TestCreateOrder(unittest.TestCase):

    @patch("client.requests.post")
    def test_create_order_success(self, mock_post):
        mock_response = MagicMock()
        mock_response.json.return_value = SAMPLE_RESPONSE
        mock_response.raise_for_status = MagicMock()
        mock_post.return_value = mock_response

        result = create_order(make_sample_order())

        self.assertEqual(result["id"], 1)
        self.assertEqual(result["status"], "pending")
        self.assertAlmostEqual(result["total_amount"], 1589.96, places=2)

    @patch("client.requests.post")
    def test_create_order_sends_correct_payload(self, mock_post):
        mock_response = MagicMock()
        mock_response.json.return_value = SAMPLE_RESPONSE
        mock_response.raise_for_status = MagicMock()
        mock_post.return_value = mock_response

        create_order(make_sample_order())

        call_kwargs = mock_post.call_args
        payload = call_kwargs.kwargs["json"]
        self.assertEqual(payload["customer_id"], 101)
        self.assertEqual(len(payload["items"]), 2)
        self.assertEqual(payload["ship_to"]["city"], "Saint Petersburg")

    @patch("client.requests.post")
    def test_create_order_raises_on_http_error(self, mock_post):
        import requests as req
        mock_response = MagicMock()
        mock_response.raise_for_status.side_effect = req.HTTPError("400 Bad Request")
        mock_post.return_value = mock_response

        with self.assertRaises(req.HTTPError):
            create_order(make_sample_order())


class TestGetOrder(unittest.TestCase):

    @patch("client.requests.get")
    def test_get_order_success(self, mock_get):
        mock_response = MagicMock()
        mock_response.json.return_value = SAMPLE_RESPONSE
        mock_response.raise_for_status = MagicMock()
        mock_get.return_value = mock_response

        result = get_order(1)

        self.assertEqual(result["id"], 1)
        self.assertEqual(result["customer_id"], 101)
        mock_get.assert_called_once_with("http://localhost:8082/orders/1", timeout=10)

    @patch("client.requests.get")
    def test_get_order_not_found(self, mock_get):
        import requests as req
        mock_response = MagicMock()
        mock_response.raise_for_status.side_effect = req.HTTPError("404 Not Found")
        mock_get.return_value = mock_response

        with self.assertRaises(req.HTTPError):
            get_order(999)


class TestListOrders(unittest.TestCase):

    @patch("client.requests.get")
    def test_list_orders_returns_list(self, mock_get):
        mock_response = MagicMock()
        mock_response.json.return_value = [SAMPLE_RESPONSE]
        mock_response.raise_for_status = MagicMock()
        mock_get.return_value = mock_response

        result = list_orders()

        self.assertIsInstance(result, list)
        self.assertEqual(len(result), 1)

    @patch("client.requests.get")
    def test_list_orders_empty(self, mock_get):
        mock_response = MagicMock()
        mock_response.json.return_value = []
        mock_response.raise_for_status = MagicMock()
        mock_get.return_value = mock_response

        result = list_orders()
        self.assertEqual(result, [])


if __name__ == "__main__":
    unittest.main()
