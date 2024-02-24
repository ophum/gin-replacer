const apiBaseURL = "%APIBASEURL%";

document.getElementById("requestButton").addEventListener("click", async () => {
    const res = await fetch(apiBaseURL+"/ping")
    document.getElementById("response").innerText = (await res.json()).message;
})