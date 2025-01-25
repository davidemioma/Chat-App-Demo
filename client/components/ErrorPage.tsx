import Link from "next/link";

type Props = {
  title: string;
  heading: string;
  subheading?: string;
  linkText?: string;
  href?: string;
};

export default function ErrorPage({
  title,
  heading,
  subheading,
  linkText,
  href,
}: Props) {
  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="text-center">
        <h1 className="mb-4 text-6xl font-bold">{title}</h1>

        <h2 className="mb-4 text-3xl font-semibold">{heading}</h2>

        {subheading && (
          <p className="mb-8 text-xl text-gray-400">{subheading}</p>
        )}

        {linkText && href && (
          <Link
            href={href}
            className="rounded bg-blue-500 px-4 py-2 font-bold text-white transition duration-300 hover:bg-blue-600"
          >
            {linkText}
          </Link>
        )}
      </div>
    </div>
  );
}
