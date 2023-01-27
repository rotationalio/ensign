import busyOtters from '/src/assets/images/busy-sea-otters.png';

export default function MaintenaceMode() {
  return (
    <section className="max-w-4xl rounded-lg border border-solid border-primary-800 text-2xl">
      <p className="mx-auto mt-8 max-w-[600px]">
        Ensign is temporarily undergoing scheduled maintenace. We&#39;ll be back online shortly. See
        our <a href="#">operating status</a> for additional information. Contact us any time with
        questions.
      </p>
      <img src={busyOtters} alt="" className="mx-auto mt-10" />
      <p className="mt-8 pb-20 text-center">
        Enjoy a cup of coffee and catch up on our latest{' '}
        <span className="font-bold text-primary">
          <a href="https://rotational.io/blog">blog posts!</a>
        </span>
      </p>
    </section>
  );
}
