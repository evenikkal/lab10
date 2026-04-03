"""
Python client for the Go JSON Exchange service (task M5).
Demonstrates sending and receiving complex nested JSON structures.
"""

import requests
import json
from dataclasses import dataclass, asdict
from typing import List
from datetime import datetime


GO_SERVICE_URL = "http://localhost:8082"


@dataclass
class Address:
    street: str
    city: str
    country: str
    zip: str


@dataclass
class OrderItem:
    product_id: int
    product_name: str
    quantity: int
    unit_price: float


@dataclass
class Order:
    customer_id: int
    items: List[OrderItem]
    ship_to: Address


def create_order(order: Order) -> dict:
    """Send a complex nested Order to the Go service."""
    payload = {
        "customer_id": order.customer_id,
        "items": [asdict(item) for item in order.items],
        "ship_to": asdict(order.ship_to),
    }
    response = requests.post(f"{GO_SERVICE_URL}/orders", json=payload, timeout=10)
    response.raise_for_status()
    return response.json()


def get_order(order_id: int) -> dict:
    """Fetch an order by ID from the Go service."""
    response = requests.get(f"{GO_SERVICE_URL}/orders/{order_id}", timeout=10)
    response.raise_for_status()
    return response.json()


def list_orders() -> List[dict]:
    """Fetch all orders."""
    response = requests.get(f"{GO_SERVICE_URL}/orders", timeout=10)
    response.raise_for_status()
    return response.json()


def print_order(order: dict) -> None:
    """Pretty-print an order."""
    print(json.dumps(order, indent=2, default=str))


if __name__ == "__main__":
    order = Order(
        customer_id=101,
        items=[
            OrderItem(product_id=1, product_name="Laptop Pro", quantity=1, unit_price=1499.99),
            OrderItem(product_id=2, product_name="USB Hub", quantity=3, unit_price=29.99),
        ],
        ship_to=Address(
            street="Lenina 5, apt 12",
            city="Saint Petersburg",
            country="Russia",
            zip="190000",
        ),
    )

    print("Creating order...")
    created = create_order(order)
    print_order(created)

    order_id = created["id"]
    print(f"\nFetching order #{order_id}...")
    fetched = get_order(order_id)
    print_order(fetched)

    print("\nListing all orders...")
    all_orders = list_orders()
    print(f"Total orders: {len(all_orders)}")
