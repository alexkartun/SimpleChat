new Vue({
	el: '#app',

	data: {
		ws: null, 				// our websocket
		newMsg: '', 			// holds new messages to be sent to the server
		chatContent: '',        // a running list of chat messages displayed on the screen
		email: null,			// email address used for grabbing an avatar
		username: null,			// our username
		joined: false			// true if email and username filled in
	},

	created: function() {
		var self = this;
		this.ws = new WebSocket('ws://' + window.location.host + '/ws');
		this.ws.addEventListener('message', function(e) {
			var msg = JSON.parse(e.data);
			self.chatContent += '<div class="chip">'
                    + '<img src="' + self.gravatarURL(msg.email) + '">' 				// Avatar
                    + msg.username
                + '</div>'
                + emojione.toImage(msg.message) + '<br/>'; 								// Parse emojis
		});
	},

	updated: function() {
		var element = document.getElementById('chat-messages');
			element.scrollTop = element.scrollHeight;									// auto scroll to the bottom
	},

	methods: {
		send: function() {
			if (this.newMsg != '') {
				this.ws.send(
					JSON.stringify({
						email: this.email,
						username: this.username,
						message: $('<p>').html(this.newMsg).text()						// strip out html
					})
				);
				this.newMsg = '';														// reset newMsg	
			}
		},

		join: function() {
			if (!this.email) {
				Materialize.toast('You must entrer an email', 2000);
				return
			}
			if (!this.username) {
				Materialize.toast('You must choose a username', 2000);
				return
			}

			this.email = $('<p>').html(this.email).text();
			this.username = $('<p>').html(this.username).text();
			this.joined = true;
		},

		gravatarURL: function(email) {
            return 'https://s.gravatar.com/avatar/' + CryptoJS.MD5(email);
        }
	}
}
)