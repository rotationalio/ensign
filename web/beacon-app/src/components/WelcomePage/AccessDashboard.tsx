import checkmark from '/src/assets/images/checkmark.png';

function AccessDashboard() {
  return (
    <div>
      {/*     Make green cirlce the background image
       */}{' '}
      <img src={checkmark} alt="" />
      <div>
        <a href="#">Access Dashboard</a>
      </div>
    </div>
  );
}

export default AccessDashboard;
