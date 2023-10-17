
const geChars = "აბგდევზთიკლმნოპჟრსტუფქღყშჩცძწჭხჯჰ" 

let search = '';
const wordInput = document.getElementById('word');
const hintsContainer = document.getElementById('hints');
const translationContainer = document.getElementById('translation');
const clear = document.getElementById('clear');
const form = document.getElementById('form1');
const tab1header = document.getElementById('tab-header1');
const tab2header = document.getElementById('tab-header2');
const theme = document.getElementById('theme');
const tab1 = document.getElementById('tab1');
const tab2 = document.getElementById('tab2');
const wordlist = document.getElementById('wordlist');

function isCharFrom(ch, charSet) {
    for (i = 0; i < charSet.length; i++) {
        if (charSet.charAt(i) == ch) {
            return true
        }
    }
    return false
} 

function isStringGe(text) {
    for (i = 0; i < text.length; i++) {
    	if (isCharFrom(text.charAt(i), geChars)) return true;
    }
    return false;
}

function handleResponse(response) {
    if (response.success && response.items != undefined) {        
        hintsContainer.innerText = '';
        translationContainer.innerText = '';
        let i = 0;

        while (i < response.items.length) {
            let hintItem = document.createElement("div");
            hintItem.innerHTML = response.items[i].ru + ': ' + response.items[i].ge;
            hintItem.setAttribute("ru", response.items[i].ru);
            hintItem.setAttribute("ge", response.items[i].ge);
            hintItem.onclick = function() { wordSelected(this) };
            hintsContainer.appendChild(hintItem);
            i++;
        }
    } else {
        console.log(response.message)
    }
}

function makeRequest(msgText) {
	if (msgText == '') return;

    var http = new XMLHttpRequest();
    http.open("GET", "/search?query=" + msgText, true);
    http.send(null);
    http.onload = function() {
    	const response = JSON.parse(http.responseText);
        handleResponse(response)
    }
}

function wordSelected(item) {
    hintsContainer.innerHTML = '';
    wordInput.value = item.getAttribute("ru");
    clear.classList.remove('d-none');    
    translationContainer.innerText = item.getAttribute("ge");
}

const inputHandler = function(e) {
  if (wordInput.value == '') {
    clear.classList.add('d-none');
  } else {
    clear.classList.remove('d-none');
  }
  const text = wordInput.value.trim().toLowerCase();
  makeRequest(text)
}

function clearFunc() {
    wordInput.value = '';
    translationContainer.innerText = '';
    hintsContainer.innerText = '';
    clear.classList.add('d-none');
}

function onSubmit() {
	alert(1);
}

function selectSearchTab() {
	tab1header.classList.add("active");
	tab2header.classList.remove("active");
	tab1.style.display = "block";
	tab2.style.display = "none";
}

function selectThemeTab() {
	tab2header.classList.add("active");
	tab1header.classList.remove("active");
	tab1.style.display = "none";
	tab2.style.display = "block";	
}

function fillTheme() {
	themes.forEach((thm) => {
		const option = document.createElement("option");
		option.setAttribute("value", thm.id);
		option.innerText = thm.name;
		theme.appendChild(option);
	});
}

function onThemeSelected() {
	wordlist.innerHTML = '';
	const themeId = theme.value;
	words.filter(o => o.themeId == themeId).sort((a, b) => a.translation.localeCompare(b.translation)).forEach((w) => {
		const tr = document.createElement("tr");
		const td1 = document.createElement("td");
		const td2 = document.createElement("td");
		td2.innerText = w.word;
		td1.innerText = w.translation;
		tr.appendChild(td1);
		tr.appendChild(td2);
		wordlist.appendChild(tr);
	});
}

wordInput.addEventListener('input', inputHandler);
clear.addEventListener('click', clearFunc);
wordInput.addEventListener("keyup", function(event) {
        if (event.keyCode == 13) {
        	wordSelected(wordInput.value);
        }
    });

tab1header.addEventListener('click', selectSearchTab);
tab2header.addEventListener('click', selectThemeTab);
fillTheme();
onThemeSelected();
theme.addEventListener("change", onThemeSelected);

