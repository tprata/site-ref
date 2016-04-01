var comteur_click_connexion = 0;
var comteur_click_inscription = 0;
var i = 0;

var url = "http://localhost:8000/api"


function removeallerrors() {
	$('#erroremail').remove()
	$('#errorlogin').remove()
	$('#errorname').remove()
	$('#errorpasswd1').remove()
	$('#errorpasswd2').remove()
	$('#errorvalue').remove()
	$('#errorloginpost').remove()
	$('#erroremailpost').remove()
	$('.successdiv').remove()
}


function isValidEmailAddress(emailAddress) {
    var pattern = /^([a-z\d!#$%&'*+\-\/=?^_`{|}~\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]+(\.[a-z\d!#$%&'*+\-\/=?^_`{|}~\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]+)*|"((([ \t]*\r\n)?[ \t]+)?([\x01-\x08\x0b\x0c\x0e-\x1f\x7f\x21\x23-\x5b\x5d-\x7e\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]|\\[\x01-\x09\x0b\x0c\x0d-\x7f\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]))*(([ \t]*\r\n)?[ \t]+)?")@(([a-z\d\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]|[a-z\d\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF][a-z\d\-._~\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]*[a-z\d\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])\.)+([a-z\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]|[a-z\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF][a-z\d\-._~\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]*[a-z\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])\.?$/i;
    return pattern.test(emailAddress);
}

function isValidLogin(login) {
    var pattern = /^[a-z0-9_-]{4,16}$/;
    return pattern.test(login);
}

function isValidName(name) {
    var pattern = /^[a-zA-Z0-9áàâäãåçéèêëíìîïñóòôöõúùûüýÿæœÁÀÂÄÃÅÇÉÈÊËÍÌÎÏÑÓÒÔÖÕÚÙÛÜÝŸÆŒ\s]{3,50}$/;
    return pattern.test(name);
}

function isValidPassword(passwd) {
	if (passwd.length > 5)
		return true
	return false
}

function signin() {
	removeallerrors();
	var login = $('input[name=id]').val()
	var password = $('input[name=password]').val()

	var post = "id=" + login + "&password=" + password

	$.post(url + '/session/create', post, function(data, textStatus, xhr) {
		if (data.Status != "OK") {
			$('<div id="erroremail" class="errorinscription"><br><span></span></div>').appendTo("body")
			$('#erroremail').css('top', '25%');
			$('#erroremail').css('height', '8%')

			if (data.List_Errors['login']) {
				$('#erroremail span').text('L\'identifiant ou email n\'existe pas.')
			} else {
				$('#erroremail span').text('L\'identifiant ou email ne correspond pas au mot de passe.')
				$('#erroremail').css('height', '10%')
			}
		} else {
			$.cookie("token_session", data.Token);
			window.location.replace("http://localhost:8000");
		}
	});
}

function signup() {
	$('#errorvalue').remove()
	$('#errorloginpost').remove()
	$('#erroremailpost').remove()

	var login = $('input[name=login]').val()
	var email = $('input[name=email]').val()
	var name = $('input[name=name]').val()
	var password = $('input[name=password1]').val()
	var password2 = $('input[name=password2]').val()

	if ($('.errorinscription').length > 0) {
		return
	}
	if (login == "" || email == "" || name == "" || password == "" || password2 == "") {
		$('<div id="errorvalue" class="errorinscription"><br><span>Veuillez completer entièrement le formulaire pour vous inscrire.</span></div>').appendTo("body")
		$('#errorvalue').css('height', '11%')
		return
	}

	var post = "email=" + email + "&password1=" + password + "&password2=" + password2 + "&login=" + login + "&name=" + name
	$.post(url + '/user/create', post, function(data, textStatus, xhr) {
		if (data.Status == "OK") {
			$('.inscriptiondiv').remove();
			$('<div class="successdiv"><h1>Votre compte à été créer avec succès.</h1></div>').appendTo('body')
			$('<div class="connexiondiv"><h2>Connexion : </h2><br><input type="text" name="id" placeholder="identifiant ou email"><br><br><input type="password" name="password" placeholder="mot de passe"><br><br><br><input id="signin" type="submit" value="Connexion"></div>').appendTo('body')
		} else {
			if (data.List_Errors['login']) {
				$('<div id="errorloginpost" class="errorinscription"><br><span>Il existe déjà un compte portant cet identifiant. Veuillez en choisir un nouveau.</span></div>').appendTo("body")
				$('#errorloginpost').css('height', '11%')
				$('#errorloginpost').css('top', '23%');
			}
			if (data.List_Errors['email']) {
				$('<div id="erroremailpost" class="errorinscription"><br><span>Cette addresse email est liée à un compte existant.</span></div>').appendTo("body")
				$('#erroremailpost').css('height', '12%')
			}
		}
	});



}


$(function() {

	$('body').on('click', '#signin', signin);
	$('body').on('click', '#signup', signup);	

	$('.sign_in_up').mouseenter(function(){
		$(this).css('background-color', 'rgba(255, 255, 255, 0.5)')		
	});
	$('.sign_in_up').mouseleave(function(){
		$(this).css('background-color', 'rgba(0, 0, 0, 0.5)')		
	});
	$('#connexion').click(function(){
		removeallerrors();
		$('.inscriptiondiv').remove();
		$('.connexiondiv').remove();
		$(this).css('color', 'white')
		$('<div class="connexiondiv"><h2>Connexion : </h2><br><input type="text" name="id" placeholder="identifiant ou email"><br><br><input type="password" name="password" placeholder="mot de passe"><br><br><br><input id="signin" type="submit" value="Connexion"></div>').appendTo('body')
	});

	$('#inscription').click(function() {
		removeallerrors();
		$('.inscriptiondiv').remove();
		$('.connexiondiv').remove();
		$(this).css('color', 'white')
		$('<div class="inscriptiondiv"><h2>Inscription : </h2><br><input type="text" name="email" placeholder="Email"><br><br><input type="text" name="login" placeholder="Identifiant"><br><br><input type="text" name="name" placeholder="Prénom"><br><br><input type="password" name="password1" placeholder="mot de passe"><br><br><input type="password" name="password2" placeholder="mot de passe de comfirmation"><br><br><br><input id="signup" type="submit" value="Inscription"></div>').appendTo('body');
});


$('body').on('change', 'input[name=email]', function(event) {
	$('#errorvalue').remove()
	value = $(this).val()
	if (!isValidEmailAddress(value)) {
		$('#erroremail').remove()
		$('<div id="erroremail" class="errorinscription"><br><span>Votre email n\'est pas au bon format.</span></div>').appendTo("body")
	} else {
		$('#erroremail').remove()
	}
});

$('body').on('change', 'input[name=login]', function(event) {
$('#errorvalue').remove()
	value = $(this).val()
	if (!isValidLogin(value)) {
		$('#errorlogin').remove()
		$('<div id="errorlogin" class="errorinscription"><br><span>Votre identifiant n\'est pas au bon format, il doit être en 4 et 16 caractères et n\'avoir que - ou _ comme caractères spéciaux.</span></div>').appendTo("body")
		$('#errorlogin').css('top', '18%');
		$('#errorlogin').css('height', '12%')
	} else {
		$('#errorlogin').remove()
	}
});

$('body').on('change', 'input[name=name]', function(event) {

	value = $(this).val()
	if (!isValidName(value)) {
		$('#errorname').remove()
		$('<div id="errorname" class="errorinscription"><br><span>Votre prénom n\'est pas au bon format.</span></div>').appendTo("body")
		$('#errorname').css('top', '31.5%');
	} else {
		$('#errorname').remove()
	}
});

$('body').on('change', 'input[name=password2]', function(event) {

	if (!isValidPassword($('input[name=password2]').val())) {
		$('#errorpasswd1').remove()
		$('<div id="errorpasswd1" class="errorinscription"><br><span>Votre mot de passe doit contenir au moins 6 caractères.</span></div>').appendTo("body")
		$('#errorpasswd1').css('top', '40%');
		$('#errorpasswd1').css('height', '9%')
	} else {
		$('#errorpasswd1').remove()
	}

	if ($('input[name=password2]').val() !== $('input[name=password1]').val()) {
		$('#errorpasswd2').remove()
		$('<div id="errorpasswd2" class="errorinscription"><br><span>Les deux mots de passes ne sont pas identiques.</span></div>').appendTo("body")
		$('#errorpasswd2').css('top', '50%');
		$('#errorpasswd2').css('height', '9%')
	} else {
		$('#errorpasswd2').remove()
	}

});



});
