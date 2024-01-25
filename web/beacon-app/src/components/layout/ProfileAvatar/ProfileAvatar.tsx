import { getInitials } from '@/utils/getInitials';

type ProfileAvatar = {
  name: string;
};

const ProfileAvatar = ({ name }: ProfileAvatar) => {
  // console.log('[] name', name);
  return (
    <div className="flex h-[28px] w-[28px] items-center justify-center rounded-full bg-[#F7F9FB] md:bg-primary">
      <span className="text-sm font-semibold md:text-white">{getInitials(name)}</span>
    </div>
  );
};

export default ProfileAvatar;
