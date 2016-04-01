var orientation = 0
var prenom = 0
var global = 0



function getCookie(name) {
  var value = "; " + document.cookie;
  var parts = value.split("; " + name + "=");
  if (parts.length == 2) return parts.pop().split(";").shift();
}

function delete_cookie(name) {
  document.cookie = name + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}

function getPathPP(images){
	var i = 0;
	while (i < images.length) {
		if (images[i].Profile == true) {
			return images[i].ImagePath
		}
		i++;
	}
}

function PutOrientation(url, data, new_val) {
    $.ajax({
        url: url,
        type: 'PUT',
        data: data,
        success: function(result) {
            if (result.Status == "OK") {
                orientation = new_val
                divli.find('select').replaceWith('<h2>'+ new_val +'</h2>')
            }
        }
    })
}

function PutPrenom(url, data, new_val) {
    $.ajax({
        url: url,
        type: 'PUT',
        data: data,
        success: function(result) {
            if (result.Status == "OK") {
                prenom = new_val
                divli.find('input').replaceWith('<h2>'+ new_val +'</h2>')
                divli.find('button').remove()
            }
        }
    })
}



var url = "http://localhost:8000/api"
var url_image = "http://localhost:8000"

var token = getCookie('token_session')
var param = "token=" +token

var pathPP

$(function() {

    global = $(this)

	$('#logout').click(function(){
		delete_cookie("token_session")
		window.location.replace("http://localhost:8000")
	});

	$('#logout').mouseenter(function(event) {
		$(this).css('background-color', 'white');
		$(this).css('color', 'rgba(0, 0,0, 0.3)')
	});

	$('#logout').mouseleave(function(event) {
		$(this).css('background-color', 'rgba(0, 0,0, 0.3)');
		$(this).css('color', 'white')
	});

	$('li').mouseenter(function(event) {
		$('<img id="modifybutton" style="position:absolute; left:190px; top:27px;" src="/assets/images/modify.png">').appendTo($(this))
	});

	$('li').mouseleave(function(event) {
		$('#modifybutton').remove()
	});

    $('#bio').mouseenter(function(event) {
        $('<img class="editionbio" id="modifybutton" src="/assets/images/modify.png">').appendTo('#bio p')
    });

    $('#bio').mouseleave(function(event) {
        $('#modifybutton').remove()
    });

    $('li').on('click', '#modifybutton', function(event) {
        divli = $(this).parent()
        $('body #modifybutton').remove()
        if (divli.attr('id') == "orientation") {
            if (orientation == "Hétéro"){
                divli.find('h2').replaceWith('<select class="selectorientation"><option>Hétéro</option><option>Homosexuel</option><option>Bisexuel</option></select>')
            } else if (orientation == "Homosexuel"){
                divli.find('h2').replaceWith('<select class="selectorientation"><option>Homosexuel</option><option>Hétéro</option><option>Bisexuel</option></select>')
            } else {
                divli.find('h2').replaceWith('<select class="selectorientation"><option>Bisexuel</option><option>Hétéro</option><option>Homosexuel</option></select>')
            }
        } else if (divli.attr('id') == "prenom") {
            divli.find('h2').replaceWith('<input id="inputprenom" class="inputprenom" type="text" name="name"><button class="inputprenom">ok</button>')
            divli.find('input').val(prenom)
        }
        });

	$('#orientation').on('change', 'select', function(){
        $('body #modifybutton').remove()
        var value = $(this).val()

        urlput = url + '/user/profile/orientation' + '?token=' + $.cookie("token_session")
        put_val = "orientation=" + value
        PutOrientation(urlput, put_val, value)
    })

    $('#prenom').on('click', 'button', function() {
        $('body #modifybutton').remove()
        var value = $('#inputprenom').val();

        var urlput = url + '/user/profile/name' + '?token=' + $.cookie("token_session")
        var put_val = "name=" + value
        console.log(urlput)
        PutPrenom(url + '/user/profile/name' + '?token=' + $.cookie("token_session"), put_val, value)
    })



	$.ajax({
        url: url + '/user/profile/me',
        type: 'GET',
        dataType: 'json',
        data: param,
        success: function(data) {
            prenom = data.Name
        	$('#prenom h2').text(data.Name)
        	$('#bio p').text(data.Bio)
        	if (data.Images.length > 0) {
        		pathPP = getPathPP(data.Images)
        		$('.profile-picture img').attr('src', url_image+pathPP)
        	}
        	if (data.Images.length > 1) {
        		total_h = 220 * (data.Images.length - 1)
        		$('<div id="othersimage"></div>').appendTo('#bio')
        		$('#othersimage').css('height', total_h + "px");
        		$.each(data.Images, function(index, value) {
        			if (value.Profile == false) {
        				pathImage = value.ImagePath
        				$('<img id="'+ "image" + value.Id +'" src="'+ url_image+pathImage+'">').appendTo('#othersimage')
        			}
        		})
        	}
        	if (data.Orientation == "") {
                orientation = "Hétéro"
        		$('#orientation h2').text('Hétéro')
        	} else {
                orientation = data.Orientation
        		$('#orientation h2').text(data.Orientation)
        	}
        	if (data.Sexe === true) {
        		$('#sexe h2').text('Femme')
        	} else {
        		$('#sexe h2').text('Homme')
        	}

        	if (data.Score < 3) {
        		$('#score h2').text('Nouveau')
        	} else {
        		$('#score h2').text(data.Score)
        	}
        }
})

})