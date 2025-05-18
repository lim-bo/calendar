#ifndef REPLAY_EVENT_H
#define REPLAY_EVENT_H

#include <QDialog>
#include "eventData.h"

namespace Ui {
class replay_event;
}

class replay_event : public QDialog
{
    Q_OBJECT

public:
    explicit replay_event(EventData originalEvent, QWidget *parent = nullptr);
    ~replay_event() override;
    EventData getRepeatedEvent() const;

private slots:
    void on_pushButton_clicked();

private:
    Ui::replay_event *ui;
    EventData originalEvent;
    EventData repeatedEvent;
};

#endif // REPLAY_EVENT_H
