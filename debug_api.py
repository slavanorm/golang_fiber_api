import sys
import requests
import os
import logging

BASE_URL = "http://localhost:8080"

verbose = False

if verbose:
    level = logging.DEBUG
else:
    level = logging.WARN

ch = logging.StreamHandler()
ch.setLevel(level)

logger = logging.Logger("", level=level)
logger.addHandler(ch)


def cleanup_db():
    db = "main.db"
    if os.path.exists(db):
        os.remove(db)
        logger.info("removed")
    else:
        logger.info("no db to remove")


def send_request(
    name, url, level=logging.WARN, method="post", should_fail=False, **kwargs
):
    kwargs["url"] = BASE_URL + "/" + url
    func = getattr(requests, method)
    try:
        r = func(**kwargs)
        rr = r.json()
    except (requests.exceptions.JSONDecodeError, requests.exceptions.ConnectionError):
        logger.error("cant load json")
        sys.exit()
    logger.log(msg=f"{method.upper()} {url}", level=level)
    logger.log(msg=f"Response: {rr}\n\n", level=level)
    failed = False
    try:
        r.raise_for_status()
    except Exception as e:
        failed = True
        if not should_fail:
            logger.exception(e, exc_info=False)
            sys.exit()
    if should_fail and not failed:
        logger.error("didnt fail!")
        sys.exit()

    return rr


try:
    LOGIN = "testuser1"
    PASSW = "testpass123!"

    # User registration and login
    user_id = (
        send_request(
            "Register",
            url="auth/register",
            method="post",
            json={
                "username": LOGIN,
                "email": "test1@example.com",
                "phone": "21234567890",
                "password": PASSW,
            },
        )
        .get("data", {})
        .get("id")
    )
    token = send_request(
        "Login", url="auth/login", data={"identity": LOGIN, "password": PASSW}
    ).get("data")
    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}
    # Estate
    estate_id = (
        send_request(
            "create_e",
            url="api/estate",
            headers=headers,
            method="post",
            json={
                "name": "Test Estate",
                "type": "flat_yearly",
                "price": 2500,
                "address": "123 Main St",
                "in_rent": True,
                "user_id": 0,
            },
        )
        .get("data", {})
        .get("id")
    )

    estate = send_request(
        "get_e", url=f"api/estate/{estate_id}", method="get", headers=headers
    )
    send_request(
        "put_e",
        url=f"api/estate/{estate_id}",
        headers=headers,
        method="put",
        json={"name": "Updated Estate", "price": 3000},
    )
    send_request(
        "del_e",
        url=f"api/estate/{estate_id}",
        headers=headers,
        method="delete",
    )
    send_request(
        "get_e",
        url=f"api/estate/{estate_id}",
        method="get",
        headers=headers,
        should_fail=True,
    )
    

    # Transaction
    tx_id = (
        send_request(
            "new_t",
            url="api/transaction",
            method="post",
            headers=headers,
            json=dict(estate_id=estate_id, amount=10),
        )
        .get("data", {})
        .get("id")
    )

    send_request(
        "new_t",
        url=f"api/transaction/{tx_id}",
        method="put",
        headers=headers,
        json=dict(amount=100.15),
    )
    send_request(
        "new_t",
        url=f"api/transaction/{tx_id}",
        method="delete",
        headers=headers,
    )

    send_request(
        "new_t",
        url=f"api/transaction/{tx_id}",
        method="get",
        headers=headers,
        should_fail=True,
    )

    # user
    send_request(
        "upd_u",
        url=f"api/user/{user_id}",
        method="put",
        headers=headers,
        json=dict(password=PASSW + "1", username="qwe"),
    )
    send_request(
        "del_u",
        url=f"api/user/{user_id}",
        method="delete",
        headers=headers,
    )
    send_request(
        "get_u",
        url=f"api/user/{user_id}",
        method="get",
        headers=headers,
        should_fail=True,
    )


finally:
    cleanup_db()
