export class DOMUtils {
	static removeEl = (identifier: string) => {
		document.querySelector(identifier)?.remove()
	};

	static getEl = (identifier: string): HTMLElement | null => {
		return document.querySelector(identifier)
	}
}


