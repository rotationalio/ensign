/* Toggles between adding and removing the "responsive" class to topnav when the user clicks on the hamburger icon */
export default function toggleResponsiveClass(e) {
  e.preventDefault();
  var x = document.getElementById("myTopnav");
  if (x.className === "topnav") {
    x.className += " responsive";
  } else {
    x.className = "topnav";
  }
}