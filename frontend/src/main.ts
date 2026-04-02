import './style.css'
type Album = {
	id: number
	title: string
	artist: string
	price: number
}

type Response = {
	message: string
	code: number
}

async function get_albums(url: string): Promise<Album[] | Response> {
	try {
		const response = await fetch(url, {
			method: "GET",
			headers: {
				"Content-Type": "application/json",
				// "X-TrxApp": "your_jwt_token",
			}
		})

		if (response.status != 200) {
			console.warn("Server did not return an OK response", response.status)
			const serverResponse = await response.json()
			return serverResponse
		}
		const data: Album[] = await response.json()
		return data

	} catch (err) {
		throw new Error("Could not get response from server")
	}

}

const URL = "http://localhost:8080/albums"
function main() {
	const btn = document.querySelector(".get-albums")
	if (!btn) {
		console.error("Could not find button with class: get-albums")
		return
	}

	btn.addEventListener("click", async (evt) => {
		const res = await get_albums(URL)
		if (Array.isArray(res)) {
			res.forEach((album) => {
				console.log(album)
			})
		} else {
			console.log("Got a message from server: ", res)
		}
	})
}

main()
