import { LoginForm } from "@/components/LoginForm"
import { IconInnerShadowTop } from "@tabler/icons-react"

export default function LoginPage() {
  return (
    <div className="grid min-h-svh lg:grid-cols-2">
      <div className="flex flex-col gap-4 p-6 md:p-10">
        <div className="flex justify-center gap-2 md:justify-start">
          <a href="#" className="flex items-center gap-2 font-medium">
            <IconInnerShadowTop className="!size-5" />
            Habit Tracker
          </a>
        </div>
        <div className="flex flex-1 items-center justify-center">
          <div className="w-full max-w-xs">
            <LoginForm />
          </div>
        </div>
      </div>
      <div className="bg-muted relative hidden lg:block">
        <img
          src="/coffee.jpg"
          alt="Image"
          className="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
        />
        <a
          href="https://unsplash.com/photos/clear-drinking-glass-with-brown-liquid-mAAcR1LR0mo"
          target="_blank"
          rel="noopener noreferrer"
          className="absolute bottom-2 right-2 bg-black/50 text-white text-xs px-2 py-1 rounded hover:bg-black/70 transition-colors"
        >
          Photo by Tavis Beck on Unsplash
        </a>
      </div>
    </div>
  )
}
