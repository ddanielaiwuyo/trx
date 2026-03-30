const serverAddr = "http://localhost:9090/upload"
async function main(data: FormData): Promise<void> {
	let form_data = data
	console.log("sending form data")
	try {
		const response = await fetch(serverAddr, {
			method: "POST",
			body: form_data,
		})

		if (response.status != 200) {
			console.log("did not get an ok status-code from server", response.status)
		}

		const body = await response.json()
		console.log("body response ->")
		console.log(body)

	} catch (err) {
		console.error(err)
	}
	return
}

async function collect_form_data() {
	let form = document.querySelector("form")
	if (!form) {
		throw new ReferenceError("Could not find form element ")
	}
	console.log("form_items", form)
	form?.addEventListener("submit", async (e) => {
		e.preventDefault()
		console.log("form -> ", form, typeof form)

		await main(new FormData(form))
	})
}

collect_form_data()
