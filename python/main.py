import asyncio
import time

import httpx
from flask import Flask, jsonify

app = Flask(__name__)


@app.route("/")
async def hello_world():
    start_time = time.time()
    demo_data = dict()

    async with httpx.AsyncClient() as client:
        while len(demo_data) < 25:
            promises = await asyncio.gather(
                *[
                    client.get("https://api.chucknorris.io/jokes/random")
                    for _ in range(0, 25 - len(demo_data))
                ]
            )

            for resp in promises:
                if resp.status_code == httpx.codes.OK:
                    response = resp.json()
                    obj = demo_data.get(response.get("id", ""), None)

                    if not obj:
                        demo_data[response.get("id", "")] = response

    app.logger.info("--- %s seconds ---" % (time.time() - start_time))

    return jsonify(
        {"length": len(demo_data), "data": [item for _, item in demo_data.items()]}
    )
