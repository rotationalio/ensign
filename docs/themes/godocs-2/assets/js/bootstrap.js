// bootstrap js components
// import Alert from "js/bootstrap/src/alert";
// import Button from "js/bootstrap/src/button";
// import Carousel from "js/bootstrap/src/carousel";
import Collapse from "js/bootstrap/src/collapse";
import Dropdown from "js/bootstrap/src/dropdown";
import Modal from "js/bootstrap/src/modal";
// import Offcanvas from "js/bootstrap/src/offcanvas";
// import Popover from "js/bootstrap/src/popover";
// import ScrollSpy from "js/bootstrap/src/scrollspy";
import Tab from "js/bootstrap/src/tab";
// import Toast from "js/bootstrap/src/toast";
// import Tooltip from "js/bootstrap/src/tooltip";


(function () {
  "use strict";

  let searchModalEl = document.getElementById('searchModal');
  let modalOpen = false;
  let searchModal = new Modal(searchModalEl, {});
  
  const params = new Proxy(new URLSearchParams(window.location.search), {
    get: (searchParams, prop) => searchParams.get(prop),
  });
  if (params.search !== null && params.search !== undefined) {
    searchModal.show();
    modalOpen = true;
  } else {
    searchModal.hide();
    modalOpen = false;
  }

  document.addEventListener('keydown', function(e) {
    if (e.key === "Escape") {
      searchModal.hide();
    } else if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
      if (!modalOpen) {
        e.preventDefault();
        searchModal.show();
        modalOpen = true;
      } else {
        e.preventDefault();
        searchModal.hide();
        modalOpen = false;
      }
    }
  });

  searchModalEl.addEventListener('hidden.bs.modal', e => {
    modalOpen = false;
  });
  
  // document.addEventListener('keydown', function(e) {
  //   let searchModal = new Modal(document.getElementById('searchModal'), {});
  //   let modalOpen = document.getElementById('searchModal').classList.contains('show');
    
  //   if (e.key === "Escape") {
  //     searchModal.hide();
  //   } else if (e.ctrlKey && e.key === 'k' || e.metaKey && e.key === 'k') {
  //     e.preventDefault();
  //     if (modalOpen) {
  //       searchModal.hide();
  //     } else {
  //       searchModal.show();
  //     }
  //   }
  // });

//   let toastElList = [].slice.call(document.querySelectorAll(".toast"));
//   let toastList = toastElList.map(function (toastEl) {
//     return new Toast(toastEl);
//   });

//   toastList.forEach(function (toast) {
//     toast.show();
//   });

//   let popoverTriggerList = [].slice.call(
//     document.querySelectorAll('[data-bs-toggle="popover"]')
//   );
//   popoverTriggerList.map(function (popoverTriggerEl) {
//     return new Popover(popoverTriggerEl);
//   });
})();