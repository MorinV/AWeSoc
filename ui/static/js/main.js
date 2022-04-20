var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

let addFriendButton = document.getElementById('addFriendButton')
addFriendButton.addEventListener("click", function () {
	addFriend(addFriendButton.dataset.friend)
})

function addFriend(friendId) {
	fetch('/addFriend', {
		headers: {"Content-Type": "application/json; charset=utf-8"},
		method: 'POST',
		body: JSON.stringify({
			friend_id: friendId
		})
	}).then(r => r.json()).then(data => {
		console.log(data);
		document.location.reload();
	})
 }
