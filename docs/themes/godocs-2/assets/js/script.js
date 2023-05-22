'use strict';

// searchToggler keyboard shortcut
const searchToggler = document.querySelectorAll('[data-search-toggler]');
searchToggler.forEach((item) => {
	let userAgentData = navigator?.userAgentData?.platform || navigator?.platform || 'unknown';

	if (userAgentData == 'macOS') {
		item.innerText = `⌘ + K`
	} else {
		item.innerText = `Ctrl + K`
	}
});


// Navbar fixed
window.onscroll = function () {
	if (document.body.scrollTop > 50 || document.documentElement.scrollTop > 50) {
		document.querySelector(".navigation").classList.add("nav-bg");
	} else {
		document.querySelector(".navigation").classList.remove("nav-bg");
	}
};

// masonry
window.onload = function () {
	let masonryWrapper = document.querySelector('.masonry-wrapper');
	// if masonryWrapper is not null, then initialize masonry
	if (masonryWrapper) {
		let masonry = new Masonry(masonryWrapper, {
			columnWidth: 1
		});
	}
};

// copy to clipboard
let blocks = document.querySelectorAll("pre");
blocks.forEach((block) => {
	if (navigator.clipboard) {
		let button = document.createElement("span");
		button.innerText = "copy";
		button.className = "copy-to-clipboard";
		block.appendChild(button);
		button.addEventListener("click", async () => {
			await copyCode(block, button);
		});
	}
});
async function copyCode(block, button) {
	let code = block.querySelector("code");
	let text = code.innerText;
	await navigator.clipboard.writeText(text);
	button.innerText = "copied";
	setTimeout(() => {
		button.innerText = "copy";
	}, 700);
}