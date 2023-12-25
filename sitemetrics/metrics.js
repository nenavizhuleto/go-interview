window.onload = (event) => {
    const site = window.location.hostname
    fetch(`http://172.16.222.31:3000/metrics?site=${site}`)
}
