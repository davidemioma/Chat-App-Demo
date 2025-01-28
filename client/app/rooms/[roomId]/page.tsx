import ChatBody from "./_components/ChatBody";
import ChatForm from "./_components/ChatForm";

export default async function RoomPage({
  params,
}: {
  params: Promise<{ roomId?: string }>;
}) {
  const { roomId } = await params;

  return (
    <div className="w-screen h-screen flex flex-col">
      <div className="px-4 md:px-6 py-2">
        <h1 className="text-xl font-bold">Room #{roomId}</h1>
      </div>

      <div className="flex-1 p-4 md:p-6 bg-gray-100 overflow-y-auto">
        <ChatBody roomId={roomId as string} />
      </div>

      <div className="w-full p-4 md:p-6">
        <ChatForm />
      </div>
    </div>
  );
}
