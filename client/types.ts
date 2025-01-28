export type UserProps = {
  id: string;
  email: string;
  username: string;
  createdAt: string;
  updatedAt: string;
};

export type RoomProps = {
  id: string;
  name: string;
};

export type MessageProps = {
  roomId: string;
  clientId: string;
  content: string;
  username: string;
  type: "recieved" | "sent";
};

export type RoomUserProps = {
  id: string;
  username: string;
};
