var Monthsyear = new Array("Janvier", "Février", "Mars", "Avril", "Mai", "Juin", "Juillet", "Août", "Septembre", "Octobre", "Novembre", "Décembre");
var server_url = "http://localhost:8000/api"


function FrenchParseDate(date) {
	var d = new Date(date);

	var tmp = d.getUTCMonth();
	var month = Monthsyear[tmp]
	var day = d.getUTCDate();
	var year = d.getFullYear();

	return day + " " + month + " " + year
}

var Post = React.createClass({

	limitSizeString: function(str, n)
	{
		if (n < str.length){
			var tmp = str.substring(0, n)
			var index = tmp.lastIndexOf(" ")
			if (index != -1) {
				tmp.substring(0, index - 1)
			}
			tmp += "..."
			return tmp
		} else {
			return str
		}
	},
	getInitialState: function() {
		var date = this.props.postInformation.CreateDate
		var content = this.props.postInformation.Content
		var title = this.props.postInformation.Title
		var important = this.props.postInformation.Important
		var id = this.props.postInformation.Id
		date = FrenchParseDate(date)
		content = this.limitSizeString(content, 125)
		title = this.limitSizeString(title, 30)

		console.log(important)
		if (important == 1) {
			important = "public"
		} else if (important == 2) {
			important = "privé"
		} else {
			important = "personnel"
		}
		return {
			date: date,
			content: content,
			title: title,
			important: important,
			Id: id
		}
	},

	render: function(){
		return (
			<div className="sousdivdescription text-center row">
			<div className="col-xs-3"><h3>{this.state.date}</h3>
			</div>
			<div className="col-xs-9">
			<span >{this.state.important}</span>
			<h4>{this.state.title}</h4>
			<quotes>{this.state.content} </quotes>
			</div>
			<p className="clear">{" ".replace(/ /g, "\u00a0")}</p>
			</div>
		)
	}

})

var PostList = React.createClass({
	render: function() {
		return (
			<div key="PostList" className="divdescription col-xs-12">
			{this.props.datapost.map( post => {
				return (
					<Post key={post.Id} postInformation={post} />
				)
			})}
			</div>
		)
	}
})

var WritePosts = React.createClass({

	componentDidMount: function() {
		var dataget = 0
		var param = 'token=31fd5905-7985-4538-973a-92017c8b4c24'
		$.ajax({
  			url: server_url + "/post",
  			type: 'GET',
  			data: param,
  			dataType: 'json',
  			success: function(data) {
  				this.setState({
  					data: data.Liste
  				})
  			}.bind(this)
		})
	},
 	 getInitialState: function() {
    	return {
    		data: []
    	};
  },

	render: function() {
		return (
			<div > <PostList datapost={this.state.data} /> </div>
		)
	}

})

ReactDOM.render(<WritePosts/>, document.getElementById('PlacePost'));