"use client"

import { useState, useEffect } from "react"
import { IconBell, IconBellOff } from "@tabler/icons-react"

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import { toast } from "sonner";
import {
  requestNotificationPermission,
  getNotificationPermission,
  isNotificationSupported,
} from "@/lib/notifications"
import { Switch } from "../ui/switch"

interface NotificationsDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function NotificationsDialog({ open, onOpenChange }: NotificationsDialogProps) {
  const [notificationsEnabled, setNotificationsEnabled] = useState(false)
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    if (open) {
      // Check current notification permission status when dialog opens
      const permission = getNotificationPermission()
      setNotificationsEnabled(permission === 'granted')
    }
  }, [open])

  const handleToggleNotifications = async (enabled: boolean) => {
    if (!isNotificationSupported()) {
      toast("Not Supported", {
        description: "Browser notifications are not supported in this browser.",
      });
      return
    }

    setIsLoading(true)

    if (enabled) {
      // Request permission to enable notifications
      const granted = await requestNotificationPermission()
      
      if (granted) {
        setNotificationsEnabled(true)
        toast("Notifications Enabled", {
          description: "You'll now receive browser notifications for habit reminders.",
        });
      } else {
        setNotificationsEnabled(false)
        toast("Permission Denied", {
          description: "Please enable notifications in your browser settings to receive habit reminders.",
        });
      }
    } else {
      // Can't programmatically disable, but we can update the UI state
      setNotificationsEnabled(false)
      toast("Notifications Disabled", {
        description: "You can re-enable notifications anytime. To fully disable, please check your browser settings.",
      });
    }

    setIsLoading(false)
  }

  const currentPermission = getNotificationPermission()
  const isSupported = isNotificationSupported()

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <IconBell className="h-5 w-5" />
            Notification Settings
          </DialogTitle>
          <DialogDescription>
            Configure browser notifications for habit reminders and updates.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          {!isSupported ? (
            <div className="flex items-center gap-2 p-4 bg-muted rounded-lg">
              <IconBellOff className="h-5 w-5 text-muted-foreground" />
              <div>
                <p className="text-sm font-medium">Not Supported</p>
                <p className="text-xs text-muted-foreground">
                  Browser notifications are not supported in this browser.
                </p>
              </div>
            </div>
          ) : (
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label htmlFor="notifications-toggle" className="text-base">
                    Browser Notifications
                  </Label>
                  <div className="text-sm text-muted-foreground">
                    Get notified when it&apos;s time to track your habits
                  </div>
                </div>
                <Switch
                  id="notifications-toggle"
                  checked={notificationsEnabled}
                  onCheckedChange={handleToggleNotifications}
                  disabled={isLoading}
                />
              </div>

              {currentPermission === 'denied' && (
                <div className="p-3 bg-orange-50 dark:bg-orange-950/50 border border-orange-200 dark:border-orange-800 rounded-md">
                  <p className="text-sm text-orange-800 dark:text-orange-200">
                    <strong>Permission Denied:</strong> To enable notifications, please allow them in your browser settings and refresh the page.
                  </p>
                </div>
              )}

              <div className="text-xs text-muted-foreground">
                <p>
                  Current status: <span className="font-medium">
                    {currentPermission === 'granted' ? 'Enabled' : 
                     currentPermission === 'denied' ? 'Blocked' : 'Not set'}
                  </span>
                </p>
              </div>
            </div>
          )}

          <div className="flex justify-end">
            <Button variant="outline" onClick={() => onOpenChange(false)}>
              Close
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
} 