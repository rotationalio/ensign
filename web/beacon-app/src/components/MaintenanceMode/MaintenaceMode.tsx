import busyOtters from '/src/components/MaintenanceMode/busy_sea_otters.png';

function MaintenaceMode() {
  return (
    <section>
      <p>Ensign is temporarily undergoing scheduled maintenace.</p>
      <p>
        We&#39;ll be back online shortly. See our <a href="#">operating status</a> for additional
        information. Contact us any time with questions.
      </p>
      <img src={busyOtters} alt="" />
      <p>
        Enjoy a cup of coffee and catch up on our latest{' '}
        <a href="https://rotational.io/blog">blog posts!</a>
      </p>
    </section>
  );
}

export default MaintenaceMode;
