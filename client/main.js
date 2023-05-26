const EventSource = require("eventsource")

function main() {
    const source = new EventSource("http://localhost:3000/sse")
    source.onmessage = (event) => {
        // console.log("OnMessage Called:")
        // console.log(event)
        // console.log(JSON.parse(event.data))
        let res = JSON.parse(event.data)
        console.log("alert status: "+res.Alert)
        console.log("notice content: "+res.Notice)
        console.log("===========================================")
        // if (event.data[0].Alert) {
        //     console.log("true !!!")
        // }
    }
}

main()