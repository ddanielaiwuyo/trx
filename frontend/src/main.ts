import './style.css'
import { DOMUtils } from "./utils"

type CanonData = {
	average_spent: number
	total_out: number
	highest_out: number
	highest_in: number
	average_in: number
	total_in: number
	out_values: number[]
	in_values: number[]

}


type Response<T> = {
	data: T
	code: number
	message: string
}


async function submitForm(form: FormData): Promise<Response<CanonData>> {
	const url = "http://localhost:8080/upload"
	try {
		const response = await fetch(url, {
			method: "POST",
			body: form,
		})

		const r: Response<CanonData> = await response.json()
		return r
	} catch (err) {
		console.error("Unexpected error occured. ", err)
		throw new Error("Could not submitForm to server")
	}
}

function main() {
	const uploadForm = document.querySelector(".upload-form") as HTMLFormElement
	if (!uploadForm) {
		console.error("Could not find .upload-form on DOM")
		return
	}

	uploadForm?.addEventListener("submit", async (evt) => {
		evt.preventDefault()
		let response = await submitForm(new FormData(uploadForm))

		if (response.message !== "Success") {
			console.warn("Server did not respond with a Success Message", response)
			return
		}

		displayCanonData(response.data)
	})


}

function displayCanonData(data: CanonData) {
	let error = document.createElement("dev-error")
	console.log("CanonData >>", data)
	DOMUtils.removeEl(".upload-form")
	let totalInEl = DOMUtils.getEl(".total-in")
	let totalOutEl = DOMUtils.getEl(".total-out")
	if (!totalInEl || !totalOutEl) {
		console.error("Could not get total-in and total-out for displayCanonData")
		error.innerText = "Could not get total-in and total-out for displayCanonData"
		return
	}

	totalInEl.innerText = `Total Money In: ${data.total_in}`
	totalOutEl.innerText = `Total Money Out: ${data.total_out}`
}



try {
	main()
} catch (err) {
	console.error("Error", err)
}
