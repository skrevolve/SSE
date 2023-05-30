const EventSource = require("eventsource");
const { stringify } = require("querystring");
const sse = new EventSource("http://localhost:3000/sse")

sse.addEventListener("notice", e => {
    const data = JSON.parse(e.data);
    console.log(`notice event: ${data.Notice}`)
})

sse.onmessage = e => {
    console.log(e)
    const data = JSON.parse(e.data);
    console.log(`normal event: ${data.Notice}`)
}

sse.onerror = e => {
    let errorMessage = String(e.message)
    if (
        errorMessage.includes("connect ECONNREFUSED") ||
        errorMessage.includes("socket hang up")
    ) {
        sse.close()
        console.log("SSE is closed")
    }
}